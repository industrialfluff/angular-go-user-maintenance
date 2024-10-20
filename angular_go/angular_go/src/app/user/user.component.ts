import { Component } from '@angular/core';
import { FormControl, Validators, FormGroup } from '@angular/forms';
import { OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService } from '../services/user.service';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  public data: any;
  userForm: FormGroup = new FormGroup({});
  isEditMode: boolean = false;  // To check if it's edit mode

  constructor(private route: ActivatedRoute, private service: UserService, private router: Router) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      // Initialize form
      this.userForm = new FormGroup({
        user_id: new FormControl({ value: '', disabled: true }), // Readonly
        user_name: new FormControl('', [Validators.required]),
        first_name: new FormControl('', [Validators.required]),
        last_name: new FormControl('', [Validators.required]),
        email: new FormControl('', [Validators.required, Validators.email]),
        user_status: new FormControl('', [Validators.required]),
        department: new FormControl('', [Validators.required])
      });

      const userId = params['user_id'];

      if (userId) {
        // If user_id exists, fetch the user data and patch the form
        this.isEditMode = true;  // Mark it as edit mode
        console.log("calling user service get");
        this.service.getUser(userId).subscribe(response => {
          this.data = response;
          this.userForm.patchValue(this.data);
        });
      } else {
        // No user_id, form remains empty for new user creation
        this.isEditMode = false;
        console.log("No user_id provided, empty form");
      }
    });
  }

  onDeleteUser(row: any) {
    console.log('Deleting user:', row);
    this.router.navigate(['/userlist']);

  }

  onSubmit(): void {
    if (this.userForm.valid) {
      // Get the full form data, including the disabled user_id
      const formData = this.userForm.getRawValue();

      console.log('Form submitted', formData);

      if (this.isEditMode) {
        // Update existing user
        console.log("Calling user service put with user_id:", formData.user_id);
        this.service.putUser(formData.user_id, formData).subscribe(response => {
          this.data = response;
          this.userForm.patchValue(this.data);  // Assuming the response contains updated data
          this.router.navigate(['/userlist']);
        });
      } else {
        // Add new user
        console.log("Calling user service post for new user");
        // Set to -1 so back end knows it's a new user
        formData.user_id = -1
        this.service.postUser(formData).subscribe(response => {
          this.data = response;
          this.userForm.patchValue(this.data);  // Assuming the response contains new user data
          this.router.navigate(['/userlist']);
        });
      }
    } else {
      console.log('Form is invalid');
    }
  }
}
