import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { UserlistComponent, User } from './userlist.component';
import { UserListService } from '../services/userlist.service';
import { UserService } from '../services/user.service';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { of, throwError } from 'rxjs';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('UserlistComponent', () => {
  let component: UserlistComponent;
  let fixture: ComponentFixture<UserlistComponent>;
  let mockUserListService: jasmine.SpyObj<UserListService>;
  let mockUserService: jasmine.SpyObj<UserService>;
  let mockRouter: jasmine.SpyObj<Router>;
  let mockDialog: jasmine.SpyObj<MatDialog>;

  const sampleUsers: User[] = [
    { user_id: 1, user_name: 'jdoe', first_name: 'John', last_name: 'Doe', email: 'jdoe@example.com', user_status: 'A', department: 'Engineering' },
    { user_id: 2, user_name: 'asmith', first_name: 'Alice', last_name: 'Smith', email: 'asmith@example.com', user_status: 'I', department: 'HR' }
  ];

  beforeEach(async () => {
    // Create spies for the services
    mockUserListService = jasmine.createSpyObj('UserListService', ['getUsers']);
    mockUserService = jasmine.createSpyObj('UserService', ['deleteUser']);
    mockRouter = jasmine.createSpyObj('Router', ['navigate']);
    mockDialog = jasmine.createSpyObj('MatDialog', ['open']);

    await TestBed.configureTestingModule({
      imports: [
        UserlistComponent, // Move UserlistComponent to imports
        MatTableModule,
        MatPaginatorModule,
        MatSortModule,
        BrowserAnimationsModule
      ],
      providers: [
        { provide: UserListService, useValue: mockUserListService },
        { provide: UserService, useValue: mockUserService },
        { provide: Router, useValue: mockRouter },
        { provide: MatDialog, useValue: mockDialog }
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserlistComponent);
    component = fixture.componentInstance;
    mockUserListService.getUsers.and.returnValue(of(sampleUsers)); // Mock the user data
    fixture.detectChanges(); // Trigger component lifecycle
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load users on initialization', () => {
    component.ngOnInit();
    expect(component.dataSource.data).toEqual(sampleUsers);
  });

  it('should navigate to edit user page', () => {
    const user = sampleUsers[0];
    component.onEditUser(user);
    expect(mockRouter.navigate).toHaveBeenCalledWith(['/user', user.user_id, 'edit']);
  });

  it('should open the dialog when deleting a user', () => {
    const user = sampleUsers[0];
    mockDialog.open.and.returnValue({
      afterClosed: () => of(true) // Simulate dialog closure with 'true' response
    } as any);

    component.onDeleteUser(user);
    expect(mockDialog.open).toHaveBeenCalled();
  });

  it('should delete a user and refresh the list on confirm', fakeAsync(() => {
    const user = sampleUsers[0];
    mockDialog.open.and.returnValue({
      afterClosed: () => of(true) // Simulate dialog closure with 'true' response
    } as any);
    mockUserService.deleteUser.and.returnValue(of({})); // Mock successful deletion

    component.onDeleteUser(user);
    tick(); // Simulate async time passage
    expect(mockUserService.deleteUser).toHaveBeenCalledWith(String(user.user_id));
    expect(mockUserListService.getUsers).toHaveBeenCalled(); // Check that loadUsers was called
  }));

  it('should log error on deletion failure', fakeAsync(() => {
    const user = sampleUsers[0];
    spyOn(console, 'error'); // Spy on console.error
    mockDialog.open.and.returnValue({
      afterClosed: () => of(true) // Simulate dialog closure with 'true' response
    } as any);
    mockUserService.deleteUser.and.returnValue(throwError(() => new Error('Deletion failed')));

    component.onDeleteUser(user);
    tick();
    expect(console.error).toHaveBeenCalledWith('Error deleting user:', jasmine.any(Error));
  }));
});
