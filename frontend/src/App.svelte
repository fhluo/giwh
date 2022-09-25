<script lang="ts">
    import {GetUIDs} from "../wailsjs/go/main/App";
    import {Get5Stars, GetSharedWishName} from "../wailsjs/go/main/App.js";
    import {api} from "./models"
    import {onMount} from "svelte";

    let uidList = [] as number[]
    let uid = 0

    onMount(async () => {
        uidList = await GetUIDs();
        if (uidList.length != 0) {
            uid = uidList[1]
        }
    })

</script>

<main>
    <div class="flex flex-row">
        <div class="flex flex-col px-1 py-3 ml-3 space-y-0.5 border rounded shadow">
            <div class="menu-item">
                UID {uid}
            </div>
            <div class="menu-item">祈愿记录</div>
            <div class="menu-item">导入导出</div>
            <div class="menu-item">设置</div>
        </div>


        <div class="flex flex-col mx-5">
            <div class="flex flex-row space-x-8 shadow border rounded-lg px-6 py-3">
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

                <div>
                    <label class="inline-block font-bold">祈愿类型：</label>
                    <select class="shadow-sm border rounded-lg w-fit px-2 py-1">
                        {#each api.SharedWishTypes as wishType}
                            {#await GetSharedWishName(wishType) then wishName}
                                <option>{wishName}</option>
                                <!--                                <div class="text-lg text-gray-900">{wishName}</div>-->
                            {/await}
                        {/each}
                    </select>
                </div>
            </div>

            <div class="space-y-3 mt-5">
                {#each api.SharedWishTypes as wishType}
                    {#await Get5Stars(uid, wishType) then items}
                        {#if items.length !== 0}
                            <div class="shadow border rounded-lg py-8 px-12 tracking-wider space-y-3 w-fit">
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

                                <div class="flex flex-row flex-wrap w-full auto-cols-auto gap-x-3 gap-y-3 auto-cols-min transition duration-200">
                                    {#each items as {id, name, pulls, lang, icon}, i}
                                        <div class="space-y-1 py-4 flex flex-col w-fit select-none bg-gray-50/75 border shadow-sm rounded hover:bg-gray-100/75 transition duration-200">
                                            <div class="px-4">
                                                <div class="w-24">
                                                    <img src="{icon}"
                                                         alt="{name}"
                                                         class="pointer-events-none rounded-full border shadow bg-amber-600">
                                                </div>
                                            </div>
                                            <div class="flex flex-col items-center justify-center space-y-0.5">
                                                <span class="font-semibold text-gray-900 tracking-wide">{name}</span>
                                                <span class="text-sm text-gray-600">{pulls}</span>
                                            </div>
                                        </div>
                                    {/each}

                                </div>

                            </div>
                        {/if}
                    {/await}
                {/each}
            </div>
        </div>

    </div>
</main>

<style>
</style>
