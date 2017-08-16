import { Component, OnInit } from '@angular/core';
import { StopwatchService } from '../services/stopwatch.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-add-group-form',
  templateUrl: './add-group-form.component.html',
  styleUrls: ['./add-group-form.component.less']
})
export class AddGroupFormComponent implements OnInit {
  public group: FormGroup;

  constructor(
    private stopwatch: StopwatchService,
    private fb: FormBuilder
  ) { }

  ngOnInit() {
    this.group = this.fb.group({
      name: ['', Validators.required]
    });
  }

  public add() {
    if (this.group.valid) {
      this.stopwatch.addGroup(this.group.value).subscribe(
        () => this.group.reset(),
        e => console.log('error', e)
      );
    }
  }
}
