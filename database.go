package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// bucket names
const (
	bucketTasks  = "tasks"
	bucketGroups = "groups"
	bucketState  = "state"
	bucketSlices = "slices"
)

// StopwatchDB is a handle for accessing a stopwatch database
type StopwatchDB struct {
	db *bolt.DB
}

// NewStopwatchDB return an initialized stopwatch db
func NewStopwatchDB() *StopwatchDB {
	return &StopwatchDB{}
}

// IsOpen returns if database is open.
func (db *StopwatchDB) IsOpen() bool {
	return db.db != nil
}

// AddTask adds a task for group, using given cost code to classify time spent
func (db *StopwatchDB) AddTask(group int, task, costcode string) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t *Task = NewTask(group, task, costcode)

	// Generate new task, return task.
	err := db.db.Update(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketTasks)).Bucket(itob(group))

		// Next task ID
		id, _ := bt.NextSequence()
		t.ID = int(id)

		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(group), itob(t.ID)}, []byte("-"))
		_, err := tx.Bucket([]byte(bucketSlices)).CreateBucketIfNotExists(sliceID)
		if err != nil {
			return err
		}

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Add new task
		return bt.Put(itob(t.ID), buf)
	})

	if err != nil {
		return nil, err
	}
	return t, nil
}

// AddGroup adds a group, using given name
func (db *StopwatchDB) AddGroup(group string) (*Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var p *Group = &Group{Name: group}

	// Generate new group and return it
	err := db.db.Update(func(tx *bolt.Tx) error {
		bp := tx.Bucket([]byte(bucketGroups))

		// get next ID
		id, _ := bp.NextSequence()
		p.ID = int(id)

		// Create bucket for the group tasks
		_, err := tx.Bucket([]byte(bucketTasks)).CreateBucketIfNotExists(itob(p.ID))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		return bp.Put(itob(p.ID), buf)
	})

	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetTask returns one task details
func (db *StopwatchDB) GetTask(group, task int) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t Task
	err := db.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketTasks)).Bucket(itob(group))
		if bt == nil {
			return errors.New("group not found")
		}

		v := bt.Get(itob(task))
		if v == nil {
			return errors.New("task not found")
		}

		return json.Unmarshal(v, &t)
	})

	return &t, err
}

// GetGroup returns one group details
func (db *StopwatchDB) GetGroup(group int) (*Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var g Group
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketGroups))
		v := b.Get(itob(group))
		if v == nil {
			return errors.New("group not found")
		}

		return json.Unmarshal(v, &g)
	})

	return &g, err
}

// StartTask marks tasks start
func (db *StopwatchDB) StartTask(group, task int) (*Task, error) {
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
	err = db.db.Update(func(tx *bolt.Tx) error {
		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(group), itob(task)}, []byte("-"))
		b := tx.Bucket([]byte(bucketSlices)).Bucket(sliceID)

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
	})

	if err != nil {
		return nil, err
	}

	t.Running = &now

	return t, db.SaveTask(t)
}

// StopTask marks task stop event
func (db *StopwatchDB) StopTask(group, task int) (*Task, error) {
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
	err = db.db.Update(func(tx *bolt.Tx) error {
		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(group), itob(task)}, []byte("-"))
		b := tx.Bucket([]byte(bucketSlices)).Bucket(sliceID)

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

		return b.Put(k, []byte(now.Format(time.RFC3339)))
	})

	if err != nil {
		return nil, err
	}

	t.Running = nil
	t.Used.Add(d)

	return t, db.SaveTask(t)
}

// SaveTask updates task value in database to the given value
func (db *StopwatchDB) SaveTask(task *Task) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketTasks)).Bucket(itob(task.GroupID))
		if b == nil {
			return errors.New("group not found")
		}

		buf, _ := json.Marshal(task)
		return b.Put(itob(task.ID), buf)
	})
}

// SaveGroup updates group value in database to the given value
func (db *StopwatchDB) SaveGroup(group *Group) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketGroups))
		buf, _ := json.Marshal(group)
		return b.Put(itob(group.ID), buf)
	})
}

// GetActiveTask returns currently active task, if one is set
func (db *StopwatchDB) GetActiveTask() (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var at ActiveTask

	err := db.db.View(func(tx *bolt.Tx) error {
		bs := tx.Bucket([]byte(bucketState))
		buf := bs.Get([]byte("activeTask"))

		if buf == nil {
			return nil
		}

		if err := json.Unmarshal(buf, &at); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
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
		b := tx.Bucket([]byte(bucketState))

		at := ActiveTask{GroupID: group, TaskID: task}
		buf, err := json.Marshal(at)
		if err != nil {
			return err
		}

		log.Printf("setactive: %+v", at)
		return b.Put([]byte("activeTask"), buf)
	})
}

// ReadGroups returns all groups
func (db *StopwatchDB) ReadGroups() ([]Group, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []Group{}

	db.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketGroups)).ForEach(func(k []byte, v []byte) error {
			var p Group
			json.Unmarshal(v, &p)
			res = append(res, p)
			return nil
		})
	})

	return res, nil
}

// ReadTasks return all tasks for a group
func (db *StopwatchDB) ReadTasks(group int) ([]*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []*Task{}

	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketTasks))

		return b.Bucket(itob(group)).ForEach(func(k []byte, v []byte) error {
			var t Task
			json.Unmarshal(v, &t)
			res = append(res, &t)
			return nil
		})
	})

	return res, nil
}

// Open opens a database and initializes it
func (db *StopwatchDB) Open(path string) error {
	var err error

	// if already open...
	if db.db != nil {
		return nil
	}

	db.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	if err = db.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketState))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil

	}); err != nil {
		return err
	}

	if err = db.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketTasks))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil

	}); err != nil {
		return err
	}

	if err = db.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketSlices))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil

	}); err != nil {
		return err
	}

	if err = db.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketGroups))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil

	}); err != nil {
		return err
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

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
