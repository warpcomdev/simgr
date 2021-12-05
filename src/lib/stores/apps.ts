import { writable } from 'svelte/store';

export class App {

    appName: string;
    status: string;
    age: number;

    constructor({ appName, status, age }: { appName: string, status: string, age: number }) {
        this.appName = appName;
        this.status = status;
        this.age = age;
    }

    age_to_interval(): string {
        let seconds = this.age / 1000;
        let d = Math.floor(seconds / (3600*24));
        let h = Math.floor(seconds % (3600*24) / 3600);
        let m = Math.floor(seconds % 3600 / 60);
        let s = Math.floor(seconds % 60);
        let parts = new Array();
        if (d > 0) {
            let dDisplay = d > 0 ? d + ' ' + (d == 1 ? "día" : "días") : "";
            parts.push(dDisplay);
        }
        if (h > 0) {
            let hDisplay = h > 0 ? h + ' ' + (h == 1 ? "hora" : "horas") : "";
            parts.push(hDisplay)
        }
        if (m > 0) {
            let mDisplay = m > 0 ? m + ' ' + (m == 1 ? "minuto" : "minutos") : "";
            parts.push(mDisplay)
        }
        if (s > 0) {
            let sDisplay = s > 0 ? s + ' ' + (s == 1 ? "segundo" : "segundos") : "";
            parts.push(sDisplay)
        }
        return parts.join(', ');
    }
    
}

export const apps = writable<App[]>(null);

export async function read_apps(username: string, password: string): Promise<App[]> {
    let response = await fetch(`siddhi-apps/statistics`, {
        headers: new Headers({
            "Authorization": `Basic ${btoa(`${username}:${password}`)}`
        })
    });
    if (response.status == 404) {
        // Siddhi devuelve un 404 cuando no hay resultados.
        // Devuelvo un array vacío para distinguir del caso en que
        // no tengo credenciales.
        return new Array();
    }
    if (response.status < 200 || response.status > 204) {
        throw `Failed to collect statistics (${response.status})`
    }
    return (await response.json()).map((item) => new App({
        appName: item.appName,
        status: item.status,
        age: Math.abs(Number(item.age || 0))
    }));
}

export async function remove_app(username: string, password: string, app: App): Promise<void> {
    let response = await fetch(`siddhi-apps/${app.appName}`, {
        method: 'DELETE',
        headers: new Headers({
            "Authorization": `Basic ${btoa(`${username}:${password}`)}`
        })
    });
    await response.text();
}

