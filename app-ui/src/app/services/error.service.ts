import { Injectable } from '@angular/core';
import {MdSnackBar, MdSnackBarConfig} from '@angular/material';

@Injectable()
export class ErrorService {

  constructor(
    private snackBar: MdSnackBar
  ) { }

  public log(e: Error, config?: MdSnackBarConfig) {
    this.snackBar.open(e.message, 'Error', config);
  }
}
