import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class RegressionService {

  private serverIP = 'http://192.168.18.97:8080'; // Replace with your server IP

  constructor(private http: HttpClient) { }

  getRegressionResults(csvURL: string): Observable<any> {
    const endpoint = `${this.serverIP}/regression?url=${encodeURIComponent(csvURL)}`;
    return this.http.get<any>(endpoint).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<any> {
    if (error.error instanceof ErrorEvent) {
      // Client-side error
      console.error('An error occurred:', error.error.message);
    } else {
      // Server-side error
      console.error(
        `Backend returned code ${error.status}, ` +
        `body was: ${error.error}`);
    }
    // Return an observable with a user-facing error message
    return throwError('Something bad happened; please try again later.');
  }
}
