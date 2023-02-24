<script lang="ts">
    import {Get5Stars, GetUIDs, GetWishName} from '../wailsjs/go/main/App.js'
    import {api} from './models'
    import {onMount} from 'svelte'
    import SideBar from './lib/SideBar.svelte'
    import Avatar from './lib/Avatar.svelte'

    let uidList = [] as number[]
    let currentUID = 0
    let currentWishType = api.SCharacterEventWish

    onMount(async () => {
        uidList = await GetUIDs()
        if (uidList.length != 0) {
            currentUID = uidList[1]
        }
    })

</script>

<main>
    <div class="flex flex-row my-3 mx-3">
        <SideBar bind:currentUID={currentUID}></SideBar>
        <div class="flex flex-col px-5 w-full h-screen">
            <div class="flex flex-row flex-wrap  items-center w-fit  mx-12">
                {#each api.SharedWishTypes as wishType}
                    {#await GetWishName(wishType) then wishName}
                        <div class="option transition duration-200 px-8 py-2 cursor-pointer hover:bg-gray-300/25 select-none leading-relaxed tracking-wider {currentWishType===wishType?'bg-gray-300/25 border-b-2 border-b-blue-500':''}"
                             on:click={()=>currentWishType=wishType}>{wishName}</div>
                    {/await}
                {/each}
            </div>

            {#await Get5Stars(currentUID, currentWishType) then items}
                {#if items.length !== 0}
                    <div class=" py-12 px-12 tracking-wider space-y-3 w-full">
                        <div class="flex flex-row flex-wrap w-full gap-x-6 gap-y-6 transition duration-200">
                            {#each items as {id, name, pulls, lang, icon}, i}
                                <Avatar icon={icon} name={name} pulls={pulls}></Avatar>
                            {/each}
                        </div>
                    </div>
                {/if}
            {/await}
        </div>
    </div>

</main>
<style>
</style>
