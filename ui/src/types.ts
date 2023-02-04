import axios, { AxiosResponse } from "axios";

export const API_URL = "http://localhost:8080";

export interface ErLog {
  id: number;
  createdAt: Date;
  updatedAt: Date;
  deletedAt?: Date;
  data: any;
}

export async function getData(route: string, body: any): Promise<ErLog[]> {
  const response: AxiosResponse<ErLog[]> = await axios.get(
    new URL(route, API_URL).href
  );

  return response.data;
}

export async function postData(route: string, body: any): Promise<ErLog[]> {
  const response: AxiosResponse<ErLog[]> = await axios.post(
    new URL(route, API_URL).href,
    body
  );

  return response.data;
}
