import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { StopwatchService } from '../../services/stopwatch.service';
import { ErrorService } from '../../services/error.service';
import { Group } from '../../model';
@Component({
  selector: 'app-group-list',
  templateUrl: './group-list.component.html',
  styleUrls: ['./group-list.component.less']
})
export class GroupListComponent implements OnInit {
  @Input() groups: Observable<Group[]>;

  constructor(
    private router: Router,
    private stopwatch: StopwatchService,
    private err: ErrorService
  ) {}

  ngOnInit() {
  }

  public goTo(group: Group) {
    this.stopwatch.selectGroup(group).subscribe(
      (g: Group) => this.router.navigate(['/group', group.id]),
      (e: Error) => this.err.log(e)
    );
  }
}
