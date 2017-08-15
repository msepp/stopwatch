package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// bucket names
const (
	bucketTasks    = "tasks"
	bucketProjects = "projects"
	bucketState    = "state"
	bucketSlices   = "slices"
)

// StopwatchDB is a handle for accessing a stopwatch database
type StopwatchDB struct {
	db *bolt.DB
}

// ActiveTask identifies currently active task
type ActiveTask struct {
	ProjectID int
	TaskID    int
}

// NewStopwatchDB return an initialized stopwatch db
func NewStopwatchDB() *StopwatchDB {
	return &StopwatchDB{}
}

// IsOpen returns if database is open.
func (db *StopwatchDB) IsOpen() bool {
	return db.db != nil
}

// AddTask adds a task for project, using given cost code to classify time spent
func (db *StopwatchDB) AddTask(project int, task, costcode string) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t *Task = NewTask(project, task, costcode)

	// Generate new task, return task.
	err := db.db.Update(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketTasks)).Bucket(itob(project))

		// Next task ID
		id, _ := bt.NextSequence()
		t.ID = int(id)

		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(project), itob(t.ID)}, []byte("-"))
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

// AddProject adds a project, using given name
func (db *StopwatchDB) AddProject(project string) (*Project, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var p *Project = &Project{Name: project}

	// Generate new project and return it
	err := db.db.Update(func(tx *bolt.Tx) error {
		bp := tx.Bucket([]byte(bucketProjects))

		// get next ID
		id, _ := bp.NextSequence()
		p.ID = int(id)

		// Create bucket for the project tasks
		_, err := tx.Bucket([]byte(bucketTasks)).CreateBucketIfNotExists(itob(p.ID))
		if err != nil {
			return err
		}

		// Add new project to projects
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
func (db *StopwatchDB) GetTask(project, task int) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var t Task
	err := db.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketTasks)).Bucket(itob(project))
		if bt == nil {
			return errors.New("project not found")
		}

		v := bt.Get(itob(task))
		if v == nil {
			return errors.New("task not found")
		}

		return json.Unmarshal(v, &t)
	})

	return &t, err
}

// StartTask marks tasks start
func (db *StopwatchDB) StartTask(project, task int) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	t, err := db.GetTask(project, task)
	if err != nil {
		return nil, err
	}

	// Find task from slices and make sure last value is an timestamp (key) without
	// end date (value)
	var now time.Time
	err = db.db.Update(func(tx *bolt.Tx) error {
		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(project), itob(task)}, []byte("-"))
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

		now = time.Now()
		return b.Put([]byte(now.Format(time.RFC3339)), []byte{})
	})

	if err != nil {
		return nil, err
	}

	t.Running = &now

	return t, db.SaveTask(t)
}

// StopTask marks task stop event
func (db *StopwatchDB) StopTask(project, task int) (*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	var d time.Duration

	t, err := db.GetTask(project, task)
	if err != nil {
		return nil, err
	}

	// Find task from slices and add end date (value) for last entry if not value
	// is yet set.
	err = db.db.Update(func(tx *bolt.Tx) error {
		// Create bucket for task slices
		sliceID := bytes.Join([][]byte{itob(project), itob(task)}, []byte("-"))
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
		now := time.Now()
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
		b := tx.Bucket([]byte(bucketTasks)).Bucket(itob(task.ProjectID))
		if b == nil {
			return errors.New("project not found")
		}

		buf, _ := json.Marshal(task)
		return b.Put(itob(task.ID), buf)
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

	if at.ProjectID == 0 && at.TaskID == 0 {
		return nil, nil
	}

	return db.GetTask(at.ProjectID, at.TaskID)
}

// SetActiveTask sets currently active task, if one is set
func (db *StopwatchDB) SetActiveTask(project, task int) error {
	if db.IsOpen() == false {
		return errors.New("database not ready")
	}

	return db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketState))

		at := ActiveTask{ProjectID: project, TaskID: task}
		buf, err := json.Marshal(at)
		if err != nil {
			return err
		}

		return b.Put([]byte("activeTask"), buf)
	})
}

// ReadProjects returns all users projects
func (db *StopwatchDB) ReadProjects() ([]Project, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []Project{}

	db.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucketProjects)).ForEach(func(k []byte, v []byte) error {
			var p Project
			json.Unmarshal(v, &p)
			res = append(res, p)
			return nil
		})
	})

	return res, nil
}

// ReadTasks return all tasks for a project
func (db *StopwatchDB) ReadTasks(project int) ([]*Task, error) {
	if db.IsOpen() == false {
		return nil, errors.New("database not ready")
	}

	res := []*Task{}

	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketTasks))

		return b.Bucket(itob(project)).ForEach(func(k []byte, v []byte) error {
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
		_, err = tx.CreateBucketIfNotExists([]byte(bucketProjects))
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
