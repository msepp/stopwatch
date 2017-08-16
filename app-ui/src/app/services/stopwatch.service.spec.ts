import { TestBed, inject } from '@angular/core/testing';

import { StopwatchService } from './stopwatch.service';

describe('StopwatchService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [StopwatchService]
    });
  });

  it('should be created', inject([StopwatchService], (service: StopwatchService) => {
    expect(service).toBeTruthy();
  }));
});
