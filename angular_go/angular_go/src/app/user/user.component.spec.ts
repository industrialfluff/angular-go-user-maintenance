import { ComponentFixture, TestBed } from '@angular/core/testing';
import { UserComponent } from './user.component';
import { ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';
import { UserService } from '../services/user.service';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('UserComponent', () => {
  let component: UserComponent;
  let fixture: ComponentFixture<UserComponent>;
  let userService: UserService;
  let router: Router;

  const mockUserService = {
    getUser: jasmine.createSpy('getUser').and.returnValue(of({
      user_id: '1',
      user_name: 'testUser',
      first_name: 'Test',
      last_name: 'User',
      email: 'test@example.com',
      user_status: 'active',
      department: 'IT'
    })),
    putUser: jasmine.createSpy('putUser').and.returnValue(of({})),
    postUser: jasmine.createSpy('postUser').and.returnValue(of({}))
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [UserComponent],
      imports: [ReactiveFormsModule, RouterTestingModule, HttpClientTestingModule],
      providers: [
        { provide: UserService, useValue: mockUserService },
        {
          provide: ActivatedRoute,
          useValue: { params: of({ user_id: '1' }) } // Mock route params
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(UserComponent);
    component = fixture.componentInstance;
    userService = TestBed.inject(UserService);
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize the form on ngOnInit', () => {
    component.ngOnInit();
    expect(component.userForm).toBeTruthy();
    expect(component.userForm.controls['user_name']).toBeDefined();
    expect(component.userForm.controls['email']).toBeDefined();
  });

  it('should load user data when user_id is present (edit mode)', () => {
    component.ngOnInit();
    expect(userService.getUser).toHaveBeenCalledWith('1');
    expect(component.isEditMode).toBeTrue();
    expect(component.userForm.value.user_name).toBe('testUser');
  });

  it('should set isEditMode to false if no user_id is present', () => {
    TestBed.overrideProvider(ActivatedRoute, { useValue: { params: of({}) } });
    component.ngOnInit();
    expect(component.isEditMode).toBeFalse();
  });

  it('should call putUser on submit in edit mode', () => {
    component.ngOnInit();
    component.userForm.controls['user_name'].setValue('updatedUser');
    component.onSubmit();

    expect(userService.putUser).toHaveBeenCalledWith('1', jasmine.objectContaining({
      user_name: 'updatedUser'
    }));
  });

  it('should call postUser on submit if in new user mode', () => {
    TestBed.overrideProvider(ActivatedRoute, { useValue: { params: of({}) } });
    component.ngOnInit();
    component.userForm.controls['user_name'].setValue('newUser');
    component.onSubmit();

    expect(userService.postUser).toHaveBeenCalledWith(jasmine.objectContaining({
      user_name: 'newUser'
    }));
  });

  it('should navigate to /userlist after deleting a user', () => {
    spyOn(router, 'navigate');
    component.onDeleteUser({ user_id: '1' });
    expect(router.navigate).toHaveBeenCalledWith(['/userlist']);
  });

  it('should log "Form is invalid" if form is invalid on submit', () => {
    spyOn(console, 'log');
    component.userForm.controls['user_name'].setValue(''); // Invalidating the form
    component.onSubmit();
    expect(console.log).toHaveBeenCalledWith('Form is invalid');
  });
});
