import { AdminUserList } from '@/components/admin/AdminUserList';

export default function AdminUsersPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Пользователи</h1>
      <AdminUserList />
    </div>
  );
}
