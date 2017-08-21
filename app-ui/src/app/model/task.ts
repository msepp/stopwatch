export class Task {
  constructor(
    public id?: number,
    public groupid?: number,
    public name?: string,
    public costcode?: string,
    public duration?: string,
    public running?: string
  ) {}
}

export function taskMatch(a: Task, b: Task): boolean {
  return a.id === b.id && a.groupid === b.groupid;
}
