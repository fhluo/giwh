<script lang="ts">
    import {Get5StarProgress, GetUIDs} from "../wailsjs/go/main/App";
    import {Get5Stars, GetPity, GetSharedWishName} from "../wailsjs/go/main/App.js";
    import {api} from "./models"
    import {onMount} from "svelte";

    let uidList = [] as number[]
    let uid = 0

    onMount(async () => {
        uidList = await GetUIDs();
        if (uidList.length != 0) {
            uid = uidList[0]
        }
    })
</script>

<main>
    <div class="flex flex-col mx-5 my-5">
        <div class="space-x-1">
            <label class="inline-block font-bold">UID:</label>
            <select class="shadow-sm border rounded-lg w-fit px-2 py-1" bind:value={uid}>
                {#await GetUIDs() then uidList}
                    {#each uidList as uid, i}
                        <option value={uid}>{uid}</option>
                    {/each}
                {/await}
            </select>
        </div>

        <div class="space-y-3 mt-5">
            {#each api.SharedWishTypes as wishType}
                <div class="shadow-sm border rounded py-3 px-5 tracking-wider space-y-3">
                    {#await GetSharedWishName(wishType) then wishName}
                        <h1 class="text-lg text-gray-900">{wishName}</h1>
                    {/await}

                    {#if wishType !== api.BeginnersWish}
                        <div class="flex flex-row items-center space-x-3">
                            <span>5 star</span>
                            {#await Get5StarProgress(uid, wishType) then progress}
                                {#await GetPity(api.Star5, wishType) then pity}
                                    <progress class="w-56" value="{progress}" max="{pity}"></progress>
                                {/await}
                            {/await}
                        </div>
                    {/if}

                    <div class="flex flex-row flex-wrap space-x-3 space-y-1">
                        {#await Get5Stars(uid, wishType) then items}
                            {#each items as {id, name, pulls}, i}
                                <div class="border shadow-sm rounded px-3 py-1 flex flex-row justify-center items-center leading-relaxed">
                                    <span class="font-bold text-yellow-500 tracking-wide">{name}</span>
                                    <span class="text-sm text-gray-300">({pulls})</span>
                                </div>
                            {/each}
                        {/await}
                    </div>

                </div>
            {/each}
        </div>
    </div>
</main>

<style>
</style>
