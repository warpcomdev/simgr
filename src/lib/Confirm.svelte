<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { App } from "$lib/stores/apps";

    // Application to show
    export let selected_app: App = null;
    // True to show the modal
    export let active = false;
    export let disabled = false;
    // Progress bar to delay this a bit.
    export let progress = 0;

    const dispatchRollback = createEventDispatcher<{rollback: App}>();
    const dispatchCommit = createEventDispatcher<{commit: App}>();

    $: {
        // Cuando el popup no esté activo, resetear disabled a false.
        disabled = disabled && active;
    }
</script>

<div class="modal" class:is-active={active}>
    <div class="modal-background"></div>
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">Confirmar borrado de aplicación</p>
            <button class="delete" on:click={() => { dispatchRollback('rollback', selected_app); }}/>
        </header>
        <section class="modal-card-body">
        <p>¿Confirma que desea eliminar la aplicación {selected_app ? selected_app.appName : ""}?</p>
        <progress class="progress is-danger" value={progress} max="100">{progress}</progress>
        </section>
        <footer class="modal-card-foot">
            <div class="buttons">
                <button class="button is-danger" disabled={disabled} on:click={() => { disabled = true; dispatchCommit('commit', selected_app); }}>Eliminar</button>
                <button class="button is-normal" on:click={() => { dispatchRollback('rollback', selected_app); }}>Cancelar</button>
            </div>
        </footer>
    </div>
</div>
