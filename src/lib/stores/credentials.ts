import { writable } from 'svelte/store';

export class Credentials {
    username: string;
    password: string;

    constructor({ username, password }: { username: string, password: string }) {
        this.username = username;
        this.password = password;
    }
}

export const credentials = writable(new Credentials({ username: "admin", password: "" }));
