import { Injectable } from '@angular/core';
import { Astilectron } from '../astilectron';
import { Observable } from 'rxjs/Observable';
import { AppConfig } from '../model';

@Injectable()
export class StopwatchService {

  constructor(
    private backend: Astilectron
  ) { }

  public config(): Observable<AppConfig> {
    return;
  }
}
