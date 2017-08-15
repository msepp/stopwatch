import { Component, OnInit } from '@angular/core';
import { Astilectron, Message } from './astilectron';
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
  public title = 'go-astilectron template project';
  public projects: Project[] = [];
  public versions: AppVersion = null;
  public newProject: FormGroup;
  public newTask: FormGroup;
  public activeTask: Task;
  public activeProject: Project;
  public tasks: Task[];

  constructor(
    private asti: Astilectron,
    private fb: FormBuilder
  ) {
    this.newProject = this.fb.group({
      name: ['', Validators.required]
    });

    this.newTask = this.fb.group({
      name: ['', Validators.required],
      costcode: ''
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
          this.projects.push(v);
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
        projectid: this.activeProject.id,
        name: this.newTask.get('name').value,
        costcode: this.newTask.get('costcode').value
      };

      this.asti.send(message.REQUEST_ADD_TASK, payload).subscribe(
        (m: Message) => {
          const v: Task = Object.assign(new Task, m.data);
          this.tasks.push(v);
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
        this.tasks = tasks;
        this.activeProject = p;
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
          t.running = nt.running;
          this.activeTask = t;
        },
        (e: Error) => {
          console.log('whoops:', e.message);
        }
      );
    } else {
      this.asti.send(message.REQUEST_STOP_TASK, payload).subscribe(
        (m: Message) => {
          const nt: Task = Object.assign(new Task, m.data);
          t.running = nt.running;
          t.duration = nt.duration;
          if (this.activeTask && this.activeTask.id === t.id && this.activeTask.projectid === t.projectid) {
            this.activeTask = null;
          }
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
        this.versions = v;

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
        m.data.projects.forEach(v => {
          this.projects.push(Object.assign(new Project, v));
        });

        if (m.data.activeTask) {
          const t: Task = Object.assign(new Task, m.data.activeTask);
          if (t.id !== 0) {
            this.activeTask = t;
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
