<script lang="ts">
    import {GetSharedWishTypes, GetUIDs} from "../wailsjs/go/main/App";
    import {GetPity, GetProgress, GetSharedWishName} from "../wailsjs/go/main/App.js";

    let uid = ""
    let rarities = ["4", "5"]
</script>

<main>
    <div class="flex flex-col mx-auto">
        {#await GetUIDs() then uidList}
            <select class="select select-bordered" bind:value={uid}>
                {#each uidList as uid, _}
                    <option>{uid}</option>
                {/each}
            </select>
        {/await}
        {#await GetSharedWishTypes() then sharedWishTypes}
            {#each sharedWishTypes as wishType, _}
                {#await GetSharedWishName(wishType) then wishName}
                    <h1 class="text-xl font-bold mt-5 ml-5">{wishName}</h1>
                {/await}
                <div class="flex flex-col">
                    {#each rarities as rarity, _}
                        <div class="ml-5">
                            <span>{rarity} star</span>
                            {#await GetProgress(uid, wishType, rarity) then progress}
                                {#await GetPity(rarity, wishType) then pity}
                                    <progress class="progress w-56" value="{progress}" max="{pity}"></progress>
                                {/await}
                            {/await}

                        </div>
                    {/each}
                </div>
                <!--            <div class="flex flex-row flex-wrap space-x-3 ml-5">-->
                <!--                {#each items as {name, pulls}, i}-->
                <!--                    <div class="font-bold bg-yellow-300 rounded-lg px-4 py-2 my-3">{name}<span-->
                <!--                            class="text-sm text-gray-500">({pulls})</span></div>-->
                <!--                {/each}-->
                <!--            </div>-->
            {/each}
        {/await}
        <!--        <button class="btn w-20 justify-self-center self-center my-3" on:click={stat}>更新</button>-->
    </div>
</main>

<style>
</style>
