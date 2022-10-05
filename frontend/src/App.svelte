<script lang="ts">
    import {GetUIDs} from "../wailsjs/go/main/App";
    import {Get5Stars, GetWishName} from "../wailsjs/go/main/App.js";
    import {api} from "./models"
    import {onMount} from "svelte";

    let uidList = [] as number[]
    let currentUID = 0
    let currentWishType = api.SCharacterEventWish

    onMount(async () => {
        uidList = await GetUIDs();
        if (uidList.length != 0) {
            currentUID = uidList[1]
        }
    })

</script>

<main>
    <div class="flex flex-col px-5 w-full h-screen">
        <div class="flex flex-row shadow border rounded-lg items-center bg-white/25">
            <label class="inline-block font-bold px-5 py-2 bg-white/75 rounded-l-lg">UID</label>
            {#await GetUIDs() then uidList}
                {#each uidList as uid, i}
                    <div class="px-5 py-2 cursor-pointer hover:bg-gray-300/50 hover:shadow hover:-top-1 select-none"
                         on:click={()=>currentUID=uid}>{uid}</div>
                {/each}
            {/await}

        </div>


        <div class="flex flex-row flex-wrap  items-center w-fit  mx-12">
            <!--                <label class="inline-block font-bold select-none px-5 py-2 bg-white/75 rounded-l-lg">祈愿类型</label>-->
            <!--                    <select class="shadow-sm border rounded-lg w-fit px-2 py-1">-->
            {#each api.SharedWishTypes as wishType}
                {#await GetWishName(wishType) then wishName}
                    <div class="option transition duration-200 px-8 py-2 cursor-pointer hover:bg-gray-300/25 select-none leading-relaxed tracking-wider {currentWishType===wishType?'bg-gray-300/25 border-b-2 border-b-blue-500':''}"
                         on:click={()=>currentWishType=wishType}>{wishName}</div>
                    <!--                                <div class="text-lg text-gray-900">{wishName}</div>-->
                {/await}
            {/each}
            <!--                    </select>-->
        </div>


        {#await Get5Stars(currentUID, currentWishType) then items}
            {#if items.length !== 0}
                <div class=" py-12 px-12 tracking-wider space-y-3 w-full">
                    <!--{#await GetSharedWishName(wishType) then wishName}-->
                    <!--    <h1 class="text-xl text-gray-900 font-semibold">{wishName}</h1>-->
                    <!--{/await}-->

                    <!--{#if wishType !== api.BeginnersWish}-->
                    <!--    <div class="flex flex-row items-center space-x-3">-->
                    <!--        <span>5 star</span>-->
                    <!--        {#await Get5StarProgress(uid, wishType) then progress}-->
                    <!--            {#await GetPity(api.Star5, wishType) then pity}-->
                    <!--                <progress class="w-56" value="{progress}" max="{pity}"></progress>-->
                    <!--            {/await}-->
                    <!--        {/await}-->
                    <!--    </div>-->
                    <!--{/if}-->

                    <div class="flex flex-row flex-wrap w-full gap-x-6 gap-y-6 transition duration-200">
                        {#each items as {id, name, pulls, lang, icon}, i}
                            <div class="cursor-pointer bg-white/50 space-y-1.5 pt-4 flex flex-col w-fit items-center select-none  border shadow-sm rounded-lg hover:bg-white/25 transition duration-200">
                                <div class="px-4 ">
                                    <div class="w-24">
                                        <img src="{icon}" alt="{name}"
                                             class="pointer-events-none rounded-full shadow-inner bg-amber-600">
                                    </div>
                                </div>
                                <div class="font-semibold text-gray-900 tracking-wide leading-relaxed">{name}</div>
                                <div class="text-sm text-gray-700 border-t py-0.5 w-full text-center bg-gray-300/25 rounded-b-lg leading-relaxed tracking-wider">{pulls}</div>
                            </div>
                        {/each}

                    </div>

                </div>
            {/if}
        {/await}
    </div>
</main>
<style>
</style>
