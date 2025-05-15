export interface User {
  id: string
  email: string
  role: string
  createdAt: string
  updatedAt: string
}

export enum UserRole {
  User = 'user',
  Admin = 'admin'
} 