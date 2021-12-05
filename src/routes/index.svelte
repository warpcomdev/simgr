<script lang="ts">
    import { credentials } from "$lib/stores/credentials";
    import { apps, read_apps, remove_app } from "$lib/stores/apps";
    import type { App } from "$lib/stores/apps";

    import Login from "$lib/Login.svelte";
    import AppTable from "$lib/AppTable.svelte";
    import Confirm from "$lib/Confirm.svelte";

    // Variables que controlan el modal de confirmación
    let selected_app: App = null;
    let active = false;
    let progress = 0;

    async function show_modal(event: CustomEvent<App>): Promise<void> {
        selected_app = event.detail;
        progress = 0;
        active = true;
    }

    async function hide_modal(): Promise<void> {
        progress = 0;
        active = false;
    }

    async function remove(event: CustomEvent<App>): Promise<void> {
        await remove_app($credentials.username, $credentials.password, event.detail);
        // Spin for a few seconds
        await fake_progress();
        $apps = await read_apps($credentials.username, $credentials.password);
        await hide_modal();
    }

    function sleep(seconds: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, seconds*1000));
    }

    async function fake_progress(): Promise<void> {
        for (let i = 0; i < 5; i++) {
            progress = 20 + i * 20;
            await sleep(1);
        }
    }
</script>

<section class="section">
    <div class="container is-max-desktop">
        <div class="box">
            <Login/>
            {#if $apps === null }
            <p>Por favor, introduzca credenciales de servidor</p>
            {:else}
                {#if $apps.length > 0}
                <AppTable on:remove={show_modal}/>
                {:else}
                <p>No hay aplicaciones</p>
                {/if}
            {/if}
        </div>
    </div>
</section>
<Confirm bind:selected_app bind:progress bind:active on:commit={remove} on:rollback={hide_modal}/>
