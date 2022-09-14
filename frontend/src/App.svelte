<script lang="ts">
    import {GetSharedWishTypes, Stat} from "../wailsjs/go/main/App.js";
    import {main} from "../wailsjs/go/models";
    import StatResult = main.StatResult;

    let results: StatResult[] = []
    let wishTypes: string[]

    function stat() {
        Stat().then(result => results = result)
    }

    function getSharedWishTypes() {
        GetSharedWishTypes().then(result => wishTypes = result)
    }

    stat()
    getSharedWishTypes()
</script>

<main>
    <div>
        <p>{wishTypes}</p>
        <h1>祈愿记录</h1>
        {#each results as {wishType, progresses}, i}
            <h1>{wishType}</h1>
            {#each progresses as {rarity, count}, i}
                <p>{rarity}: {count}</p>
            {/each}
        {/each}
        <button on:click={stat}>更新</button>
    </div>
</main>

<style>
</style>
