import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'taskDuration'
})
export class TaskDurationPipe implements PipeTransform {

  transform(value: string, args?: any): string {
    return value.replace(/([0-9])\.[0-9]+s$/, '$1s');
  }

}
