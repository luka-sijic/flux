export function getUser() {
    const token = localStorage.getItem('token');
    if (!token) {
        return null;
    }
}