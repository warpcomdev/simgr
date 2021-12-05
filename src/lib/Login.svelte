<script lang="ts">
    import { credentials } from "$lib/stores/credentials";
    import { apps, read_apps } from "$lib/stores/apps";

    async function reload() {
        $apps = await read_apps($credentials.username, $credentials.password);
	}
</script>

<form class="block">
    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label" for="username">Username</label>
        </div>
        <div class="field-body">
            <div class="field">
                <div class="control is-expanded">
                    <input class="input" type="text" name="username" id="username" placeholder="admin" bind:value={$credentials.username}/>
                </div>
            </div>
        </div>
    </div>
    <div class="field is-horizontal">
        <div class="field-label is-normal">
            <label class="label" for="password">Password</label>
        </div>
        <div class="field-body">
            <div class="field">
                <div class="control is-expanded">
                    <input class="input" type="password" name="password" id="password" bind:value={$credentials.password}/>
                </div>
            </div>
        </div>
    </div>
    <div class="field is-horizontal">
        <div class="field-label is-normal">
        </div>
        <div class="field-body">
            <div class="field">
                <div class="control">
                    <button class="button is-primary" on:click|preventDefault={reload}>Actualizar datos</button>
                </div>
            </div>
        </div>
    </div>
</form>