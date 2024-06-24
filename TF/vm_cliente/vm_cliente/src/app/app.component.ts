import { Component } from '@angular/core';
import {FormsModule} from '@angular/forms';
import {RegressionService } from './regression.service'

@Component({
  selector: 'app-root',
  standalone: true,
  template: `
    <form (ngSubmit)="submitForm()">
      <label>
        Insert .csv raw URL:
        <input id="url" type="text" [(ngModel)]="csv_url" name="csvUrl" />
      </label>
      <button type="submit" [disabled]="!csv_url">Submit</button>

      <div *ngIf="showResults">
        <p>Intercept: {{ regressionResult?.intercept }}</p>
        <p>Slope: {{ regressionResult?.slope }}</p>
      </div>
    </form>
  `,
  styleUrls: ['./app.component.css'],  // Adjust path based on your file structure
  imports: [FormsModule],
})
export class AppComponent {
  title = 'client-app';
  csv_url = '';
  showResults = false;
  regressionResult: any = null;

  constructor(private regressionService: RegressionService) {}

  submitForm() {
    this.regressionService.getRegressionResults(this.csv_url).subscribe(
      (result) => {
        this.regressionResult = result;
        this.showResults = true;
      },
      (error) => {
        console.error('Error fetching regression results:', error);
        // Handle error as needed
      }
    );
  }
}
