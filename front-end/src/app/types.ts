export type Payload = {
  status: string;
  data: string;
}

export type TestPayload = {
  [key:string]: Payload;
}