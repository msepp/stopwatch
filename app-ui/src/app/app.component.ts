import { Component, OnInit } from '@angular/core';
import { Astilectron, Message } from './astilectron';
import { Store } from '@ngrx/store';
import { AppState } from './model/app-state';
import { Observable } from 'rxjs/Observable';

import * as VersionActions from './store/actions/version.actions';
import * as ProjectsActions from './store/actions/projects.actions';
import * as ActiveTaskActions from './store/actions/active-task.actions';
import * as ProjectTasksActions from './store/actions/project-tasks.actions';
import * as SelectedProjectActions from './store/actions/selected-project.actions';

import * as message from './astilectron/message';
import { FormControl, FormGroup, FormBuilder, Validators } from '@angular/forms';

import { AppVersion } from './model/app-version';
import { Project } from './model/project';
import { Task } from './model/task';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent implements OnInit {
  public title = 'Stopwatch';
  public projects: Observable<Project[]>;
  public versions: Observable<AppVersion>;
  public activeTask: Observable<Task>;
  public selectedProject: Observable<Project>;
  public tasks: Observable<Task[]>;
  public newProject: FormGroup;
  public newTask: FormGroup;

  public projectID: number;
  public taskID: number;

  constructor(
    private asti: Astilectron,
    private fb: FormBuilder,
    private store: Store<AppState>
  ) {
    this.newProject = this.fb.group({
      name: ['', Validators.required]
    });

    this.newTask = this.fb.group({
      name: ['', Validators.required],
      costcode: ''
    });

    this.versions = this.store.select('version');
    this.projects = this.store.select('projects');
    this.tasks = this.store.select('projectTasks');
    this.activeTask = this.store.select('activeTask');
    this.selectedProject = this.store.select('selectedProject');

    this.selectedProject.subscribe((v: Project) => {
      if (v) {
        this.projectID = v.id;
      }
    });

    this.activeTask.subscribe((v: Task) => {
      if (v) {
        console.log('active task changed to', v);
        this.taskID = v.id;
      }
    });
  }

  public ngOnInit() {
    this.asti.isReady.filter(v => v === true).take(1).subscribe(() => this.getVersions());
  }

  public addProject() {
    if (this.newProject.valid) {
      this.asti.send(message.REQUEST_ADD_PROJECT, this.newProject.value).subscribe(
        (m: Message) => {
          const v: Project = Object.assign(new Project, m.data);
          this.store.dispatch(new ProjectsActions.Add(v));
          this.newProject.reset();
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  public addTask() {
    if (this.newTask.valid) {
      const payload = {
        projectid: this.projectID,
        name: this.newTask.get('name').value,
        costcode: this.newTask.get('costcode').value
      };

      this.asti.send(message.REQUEST_ADD_TASK, payload).subscribe(
        (m: Message) => {
          const v: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new ProjectTasksActions.Add(v));
          this.newTask.reset();
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  public selectProject(p: Project) {
    this.asti.send(message.REQUEST_GET_PROJECT_TASKS, {projectid: p.id}).subscribe(
      (m: Message) => {
        const tasks: Task[] = [];
        m.data.forEach(v => {
          console.log(v);
          const t: Task = Object.assign(new Task, v);
          tasks.push(t);
        });

        this.store.dispatch(new ProjectTasksActions.Set(tasks));
        this.store.dispatch(new SelectedProjectActions.Set(p));
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  public toggleTask(t: Task) {
    const payload = {
      projectid: t.projectid,
      id: t.id
    };

    if ((!!t.running) === false) {
      this.asti.send(message.REQUEST_START_TASK, payload).subscribe(
        (m: Message) => {
          const nt: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new ProjectTasksActions.Update(nt));
          this.store.dispatch(new ActiveTaskActions.Set(nt));
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    } else {
      this.asti.send(message.REQUEST_STOP_TASK, payload).subscribe(
        (m: Message) => {
          const nt: Task = Object.assign(new Task, m.data);
          this.store.dispatch(new ProjectTasksActions.Update(nt));
          this.store.dispatch(new ActiveTaskActions.Stop(nt));
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    }
  }

  private getVersions() {
    this.asti.send(message.REQUEST_APP_VERSIONS, null).subscribe(
      (m: Message) => {
        const v: AppVersion = Object.assign(new AppVersion, m.data);
        this.store.dispatch(new VersionActions.Set(v));

        this.openDatabase();
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  private openDatabase() {
    this.asti.send(message.REQUEST_OPEN_DATABASE, null).subscribe(
      () => this.getProjects(),
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }

  private getProjects() {
    this.asti.send(message.REQUEST_PROJECTS, null).subscribe(
      (m: Message) => {
        console.log(m.data);
        const all: Project[] = [];
        m.data.projects.forEach(v => {
          all.push(Object.assign(new Project, v));
        });

        this.store.dispatch(new ProjectsActions.Set(all));

        if (m.data.activeTask) {
          const t: Task = Object.assign(new Task, m.data.activeTask);
          if (t.id !== 0) {
            console.log('initial active task:', t);
            this.store.dispatch(new ActiveTaskActions.Set(t));
          }
        } else {
          this.activeTask = null;
        }
      },
      (e: Error) => {
        console.log('whoops:', e.message);
      }
    );
  }
}
