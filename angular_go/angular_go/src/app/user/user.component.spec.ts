import { ComponentFixture, TestBed } from '@angular/core/testing';
import { UserComponent } from './user.component';
import { ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';
import { UserService } from '../services/user.service';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { MatToolbarModule } from '@angular/material/toolbar'; // Import MatToolbarModule

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
      imports: [
        ReactiveFormsModule,
        RouterTestingModule,
        HttpClientTestingModule,
        MatToolbarModule // Add MatToolbarModule here
      ],
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

  // Other test cases here...
});
