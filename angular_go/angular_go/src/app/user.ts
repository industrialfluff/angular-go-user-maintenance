export interface User {
  user_id: number;
  user_name: string;  //varchar(50)
  first_name: string; //varchar(255)
  last_name: string;  //varchar(255)
  email: string;      //varchar(255)
  user_status: string;//varchar(1)
  department: string; //varchar(255) NULL
}
