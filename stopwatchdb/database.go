package stopwatchdb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	model "github.com/msepp/stopwatch/stopwatchmodel"
)

// Bucket names
const (
	BucketTasks   = "tasks"
	BucketGroups  = "groups"
	BucketState   = "state"
	BucketSlices  = "slices"
	BucketHistory = "history"
)

// StopwatchDB is a handle for accessing a stopwatch database
type StopwatchDB struct {
	db *bolt.DB
}

// New return an initialized stopwatch db
func New() *StopwatchDB {
	return &StopwatchDB{}
}

// IsOpen returns if database is open.
func (db *StopwatchDB) IsOpen() bool {
	return db.db != nil
}

// AddTask adds a task for group, using given cost code to classify time spent
func (db *StopwatchDB) AddTask(group int, task, costcode string) (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t *model.Task = model.NewTask(group, task, costcode)

	// Generate new task, return task.
	if err := db.db.Update(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(BucketTasks)).Bucket(Itob(group))

		// Next task ID
		id, _ := bt.NextSequence()
		t.ID = int(id)

		// Create Bucket for task slices
		sliceID := bytes.Join([][]byte{Itob(group), Itob(t.ID)}, []byte("-"))
		_, err := tx.Bucket([]byte(BucketSlices)).CreateBucketIfNotExists(sliceID)
		if err != nil {
			return err
		}

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Add new task
		return bt.Put(Itob(t.ID), buf)
	}); err != nil {
		return nil, err
	}

	return t, nil
}

// AddGroup adds a group, using given name
func (db *StopwatchDB) AddGroup(group string) (*model.Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var p *model.Group = &model.Group{Name: group}

	// Generate new group and return it
	if err := db.db.Update(func(tx *bolt.Tx) error {
		bp := tx.Bucket([]byte(BucketGroups))

		// get next ID
		id, _ := bp.NextSequence()
		p.ID = int(id)

		// Create Bucket for the group tasks
		_, err := tx.Bucket([]byte(BucketTasks)).CreateBucketIfNotExists(Itob(p.ID))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return bp.Put(Itob(p.ID), buf)
	}); err != nil {
		return nil, err
	}

	return p, nil
}

// GetTask returns one task details
func (db *StopwatchDB) GetTask(group, task int) (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t model.Task
	err := db.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(BucketTasks)).Bucket(Itob(group))
		if bt == nil {
			return errors.New("group not found")
		}

		v := bt.Get(Itob(task))
		if v == nil {
			return errors.New("task not found")
		}

		return json.Unmarshal(v, &t)
	})

	return &t, err
}

// GetGroup returns one group details
func (db *StopwatchDB) GetGroup(group int) (*model.Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var g model.Group
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketGroups))
		v := b.Get(Itob(group))
		if v == nil {
			return errors.New("group not found")
		}

		return json.Unmarshal(v, &g)
	})

	return &g, err
}

// StartTask marks tasks start
func (db *StopwatchDB) StartTask(group, task int) (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	t, err := db.GetTask(group, task)
	if err != nil {
		return nil, err
	}

	// Find task from slices and make sure last value is an timestamp (key) without
	// end date (value)
	var now time.Time
	if err = db.db.Update(func(tx *bolt.Tx) error {
		// Create Bucket for task slices
		sliceID := bytes.Join([][]byte{Itob(group), Itob(task)}, []byte("-"))
		b := tx.Bucket([]byte(BucketSlices)).Bucket(sliceID)

		if b == nil {
			return fmt.Errorf("task not found")
		}

		// get last value, if its value is empty, the task is already running
		c := b.Cursor()
		k, v := c.Last()
		if k != nil && len(v) == 0 {
			now, _ = time.Parse(time.RFC3339, string(k))
			return nil
		}

		now = time.Now().UTC()
		return b.Put([]byte(now.Format(time.RFC3339)), []byte{})
	}); err != nil {
		return nil, err
	}

	t.Running = &now
	return t, db.SaveTask(t)
}

// StopTask marks task stop event
func (db *StopwatchDB) StopTask(group, task int) (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var d time.Duration

	t, err := db.GetTask(group, task)
	if err != nil {
		return nil, err
	}

	// Find task from slices and add end date (value) for last entry if not value
	// is yet set.
	if err = db.db.Update(func(tx *bolt.Tx) error {
		// Create Bucket for task slices
		sliceID := bytes.Join([][]byte{Itob(group), Itob(task)}, []byte("-"))
		b := tx.Bucket([]byte(BucketSlices)).Bucket(sliceID)

		if b == nil {
			return fmt.Errorf("task not found")
		}

		// get last value, if its value is empty, the task is already running
		c := b.Cursor()
		k, v := c.Last()
		if k != nil && len(v) != 0 {
			return fmt.Errorf("task already stopped")
		}

		start, _ := time.Parse(time.RFC3339, string(k))
		now := time.Now().UTC()
		d = now.Sub(start)

		// If current slice and end point are on separate dates, we split into extra
		// slices to avoid having slices that span multiple days.
		for start.Year() != now.Year() || start.Month() != now.Month() || start.Day() != now.Day() {
			// End at next day...
			end := time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, time.UTC)
			// Minus 1 second, so last second of starting date.
			end = end.Add(time.Second * -1)

			if err := b.Put(k, []byte(end.Format(time.RFC3339))); err != nil {
				return err
			}

			// Next start is at the start of next day
			start = end.Add(time.Second)
			k = []byte(start.Format(time.RFC3339))
		}

		return b.Put(k, []byte(now.Format(time.RFC3339)))
	}); err != nil {
		return nil, err
	}

	t.Running = nil
	t.Used.Add(d)

	return t, db.SaveTask(t)
}

// SaveTask updates task value in database to the given value
func (db *StopwatchDB) SaveTask(task *model.Task) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketTasks)).Bucket(Itob(task.GroupID))
		if b == nil {
			return errors.New("group not found")
		}

		buf, _ := json.Marshal(task)
		return b.Put(Itob(task.ID), buf)
	})
}

// SaveGroup updates group value in database to the given value
func (db *StopwatchDB) SaveGroup(group *model.Group) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketGroups))
		buf, _ := json.Marshal(group)
		return b.Put(Itob(group.ID), buf)
	})
}

// GetActiveTask returns currently active task, if one is set
func (db *StopwatchDB) GetActiveTask() (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var at model.ActiveTask

	if err := db.db.View(func(tx *bolt.Tx) error {
		bs := tx.Bucket([]byte(BucketState))
		buf := bs.Get([]byte("activeTask"))

		if buf == nil {
			return nil
		}

		if err := json.Unmarshal(buf, &at); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	log.Printf("getactive: %+v", at)
	if at.GroupID == 0 && at.TaskID == 0 {
		return nil, nil
	}

	return db.GetTask(at.GroupID, at.TaskID)
}

// SetActiveTask sets currently active task, if one is set
func (db *StopwatchDB) SetActiveTask(group, task int) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketState))

		at := model.ActiveTask{GroupID: group, TaskID: task}
		buf, err := json.Marshal(at)
		if err != nil {
			return err
		}

		log.Printf("setactive: %+v", at)
		return b.Put([]byte("activeTask"), buf)
	})
}

// ReadGroups returns all groups
func (db *StopwatchDB) ReadGroups() ([]model.Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []model.Group{}

	db.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(BucketGroups)).ForEach(func(k []byte, v []byte) error {
			var p model.Group
			json.Unmarshal(v, &p)
			res = append(res, p)
			return nil
		})
	})

	return res, nil
}

// ReadTask return a task
func (db *StopwatchDB) ReadTask(group, task int) (*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var res model.Task

	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketTasks))
		bg := b.Bucket(Itob(group))
		if bg == nil {
			return errors.New("group not found")
		}

		buf := bg.Get(Itob(task))
		if buf == nil {
			return errors.New("task not found")
		}

		return json.Unmarshal(buf, &res)
	}); err != nil {
		return nil, err
	}

	return &res, nil
}

// ReadTasks return all tasks for a group
func (db *StopwatchDB) ReadTasks(group int) ([]*model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []*model.Task{}

	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketTasks))
		bg := b.Bucket(Itob(group))
		if bg == nil {
			return errors.New("group not found")
		}

		return bg.ForEach(func(k []byte, v []byte) error {
			var t model.Task
			json.Unmarshal(v, &t)
			res = append(res, &t)
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return res, nil
}

// SaveHistory updates task history to given value
func (db *StopwatchDB) SaveHistory(history []model.HistoryTask) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketHistory))
		buf, err := json.Marshal(history)
		if err != nil {
			return err
		}

		log.Printf("writing usage: %s", string(buf))
		return b.Put([]byte("usage"), buf)
	})
}

// ReadHistory returns last known task usage history
func (db *StopwatchDB) ReadHistory() ([]model.Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var history []model.HistoryTask
	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketHistory))
		buf := b.Get([]byte("usage"))
		if buf != nil {
			return json.Unmarshal(buf, &history)
		}

		// No history yet
		history = []model.HistoryTask{}
		return nil
	}); err != nil {
		return nil, err
	}

	// Retrieve tasks for history entries
	res := []model.Task{}
	for _, ht := range history {
		t, err := db.GetTask(ht.GroupID, ht.ID)
		if err == nil {
			res = append(res, *t)
		}
	}

	return res, nil
}

// Open opens a database and initializes it
func (db *StopwatchDB) Open(path string) error {
	var err error

	// Don't reopen.
	if db.db != nil {
		return nil
	}

	// Open and create if missing.
	if db.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		return err
	}

	// Then create missing buckets.
	Buckets := []string{
		BucketState,
		BucketTasks,
		BucketSlices,
		BucketGroups,
		BucketHistory,
	}
	for _, Bucket := range Buckets {
		if err = db.db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(Bucket))
			if err != nil {
				return fmt.Errorf("create Bucket '%s' failed: %s", Bucket, err)
			}
			return nil

		}); err != nil {
			return err
		}
	}

	return nil
}

// Close closes current DB handle
func (db *StopwatchDB) Close() error {
	var err error = nil

	if db.db != nil {
		err = db.db.Close()
		db.db = nil
	}

	return err
}

// Itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
