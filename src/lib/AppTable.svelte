<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { apps } from '$lib/stores/apps';
    import type { App } from '$lib/stores/apps';

    const dispatchRemove = createEventDispatcher<{remove: App}>();
</script>

<table class="table is-fullwidth">
    <thead>
        <tr>
            <th>appName</th>
            <th>status</th>
            <th>age</th>
        </tr>
    </thead>
    <tbody>
        {#each $apps as app}
        <tr>
            <td><span class="tag is-medium" class:is-success={app.status=="active"} class:is-warning={app.status!="active"}>
                {app.appName}
                <button class="delete" on:click={() => { dispatchRemove('remove', app); }}/>
            </span></td>
            <td>{app.status}</td>
            <td>{app.age_to_interval()}</td>
        </tr>
        {/each}
    </tbody>
</table>
