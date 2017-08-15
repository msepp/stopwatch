export class Task {
  constructor(
    public id?: number,
    public projectid?: number,
    public name?: string,
    public costcode?: string,
    public duration?: string,
    public running?: number
  ) {}
}
