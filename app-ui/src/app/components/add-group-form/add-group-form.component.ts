import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormGroupDirective } from '@angular/forms';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';

@Component({
  selector: 'app-add-group-form',
  templateUrl: './add-group-form.component.html',
  styleUrls: ['./add-group-form.component.less']
})
export class AddGroupFormComponent implements OnInit {
  public group: FormGroup;
  @ViewChild(FormGroupDirective) groupForm;

  constructor(
    private stopwatch: StopwatchService,
    private fb: FormBuilder,
    private err: ErrorService
  ) { }

  ngOnInit() {
    this.group = this.fb.group({
      name: ['', Validators.required]
    });
  }

  public add() {
    if (this.group.valid) {
      this.stopwatch.addGroup(this.group.value).subscribe(
        () => {
          console.log('added group');
          this.group.reset({name: ''});
          this.groupForm.resetForm();
        },
        (e: Error) => this.err.log(e)
      );
    }
  }
}
