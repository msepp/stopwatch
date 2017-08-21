import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupUsageComponent } from './group-usage.component';

describe('GroupUsageComponent', () => {
  let component: GroupUsageComponent;
  let fixture: ComponentFixture<GroupUsageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GroupUsageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GroupUsageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
