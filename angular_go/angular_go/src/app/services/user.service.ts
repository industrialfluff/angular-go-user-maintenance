import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private apiUrl = 'http://localhost:8080/users';

  constructor(private http: HttpClient) {}

  // Get a user by ID
  getUser(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/${userId}`);
  }

  // Create a new user (POST)
  postUser(userData: any): Observable<any> {
    return this.http.post(this.apiUrl, userData);
  }

  // Update a user (PUT)
  putUser(userId: string, userData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/${userId}`, userData);
  }

  // Partially update a user (PATCH)
  patchUser(userId: string, partialData: any): Observable<any> {
    return this.http.patch(`${this.apiUrl}/${userId}`, partialData);
  }

  // Delete a user (DELETE)
  deleteUser(userId: string): Observable<any> {
    console.log("deleting the user ", userId);
    return this.http.delete(`${this.apiUrl}/${userId}`);
  }
}
