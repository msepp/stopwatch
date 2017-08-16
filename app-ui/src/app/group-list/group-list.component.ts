import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { Group } from '../model';

@Component({
  selector: 'app-group-list',
  templateUrl: './group-list.component.html',
  styleUrls: ['./group-list.component.less']
})
export class GroupListComponent implements OnInit {
  @Input() groups: Observable<Group[]>;
  constructor(
    private router: Router
  ) {}

  ngOnInit() {
  }

  public goTo(group: Group) {
    this.router.navigate(['/group', group.id]);
  }
}
