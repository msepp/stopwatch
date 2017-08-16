import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs/Observable';
import { AppState, Group } from '../model';

@Component({
  selector: 'app-groups-home',
  templateUrl: './groups-home.component.html',
  styleUrls: ['./groups-home.component.less']
})
export class GroupsHomeComponent implements OnInit {
  public groups: Observable<Group[]>;
  constructor(
    private store: Store<AppState>
  ) {
    this.groups = this.store.select('groups');
   }

  ngOnInit() {
  }

}
