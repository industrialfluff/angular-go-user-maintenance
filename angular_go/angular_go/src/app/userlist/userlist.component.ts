import { Component, AfterViewInit, ViewChild } from '@angular/core';
import { UserListService } from '../services/userlist.service';
import { UserService } from '../services/user.service';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { NgIf, NgFor } from '@angular/common';
import { Router } from '@angular/router';
import { MatButton } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { ConfirmDialogComponent } from '../confirm-dialog/confirm-dialog.component';


@Component({
  selector: 'app-root',
  templateUrl: './userlist.component.html',
  styleUrls: ['./userlist.component.css'],
  standalone: true,
  imports: [MatSortModule, MatTableModule, MatPaginatorModule, NgIf, NgFor],
})
export class UserlistComponent implements AfterViewInit {

  columns = [
    {
      columnDef: 'user_id',
      header: 'User Id',
      cell: (element: User) => `${element.user_id}`,
    },
    {
      columnDef: 'user_name',
      header: 'User Name',
      cell: (element: User) => `${element.user_name}`,
    },
    {
      columnDef: 'first_name',
      header: 'First Name',
      cell: (element: User) => `${element.first_name}`,
    },
    {
      columnDef: 'last_name',
      header: 'Last Name',
      cell: (element: User) => `${element.last_name}`,
    },
    {
      columnDef: 'email',
      header: 'Email',
      cell: (element: User) => `${element.email}`,
    },
    {
      columnDef: 'department',
      header: 'Department',
      cell: (element: User) => `${element.department}`,
    },
    {
      columnDef: 'user_status',
      header: 'User Status',
      cell: (element: User) => `${element.user_status}`,
    },
  ];

  //displayedColumns = this.columns.map(c => c.columnDef);
  displayedColumns: string[] = ['user_id', 'user_name', 'first_name', 'last_name', 'email', 'department', 'user_status', 'actions'];
  dataSource = new MatTableDataSource<User>();

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;  // Use MatSort instead of MatSortModule

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;  // Correctly assign MatSort
  }

  public data: any;

  constructor(private listservice: UserListService, private userservice: UserService, private router: Router, private dialog: MatDialog) {
    console.log("calling userlist service");
  }

  ngOnInit() {
    this.loadUsers();
  }
  loadUsers() {
    this.listservice.getUsers()
    .subscribe(response => {
      this.data = response;
      this.dataSource.data = this.data;
    });
  }
  onEditUser(row: any): void {
    console.log('Edit user:', row);
    this.router.navigate(['/user', row.user_id, 'edit']);
  }

  onAddUser() {
    this.router.navigate(['/user/new']);
  }
  /*
  onDeleteUser(row: any, callback: (error?: any) => void): void {
    console.log("Deleting user with ID:", row.user_id);
    this.userservice.deleteUser(row.user_id).subscribe(
      () => {
        console.log('User deleted successfully');
        callback(); // Invoke the callback without an error
      },
      (error) => {
        console.error('Error deleting user:', error);
        callback(error); // Invoke the callback with an error
      }
    );
  }
    */
  onDeleteUser(row: any): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '300px',
      data: { message: 'Are you sure you want to delete this user?' }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Proceed with deletion only if confirmed
        this.userservice.deleteUser(row.user_id).subscribe(
          () => {
            console.log('User deleted successfully');
            this.loadUsers(); // Refresh the user list
          },
          (error) => {
            console.error('Error deleting user:', error);
            // Optionally handle the error, e.g., display a notification
          }
        );
      }
    });
  }
  onDeleteCallback = (error?: any): void => {
    if (error) {
      // Handle the error
      console.error('Error during deletion:', error);
      // Optionally, display a notification to the user
    } else {
      // Handle the success
      console.log('User deleted successfully');
      // Refresh the user list
      this.loadUsers();
    }
  };


}

export interface User {
  user_id: number;
  user_name: string;  // varchar(50)
  first_name: string; // varchar(255)
  last_name: string;  // varchar(255)
  email: string;      // varchar(255)
  user_status: string;// varchar(1)
  department: string; // varchar(255) NULL
}
