<script lang="ts">
    import {GetUIDs} from '../../wailsjs/go/main/App.js'

    export let currentUID: string

</script>

<div>
    <div class="flex flex-row shadow ring-2 ring-cyan-500 rounded items-center bg-white/25 text-sm">
        <label class="inline-block font-bold px-4 py-1 border-r bg-white/75 rounded-l-lg select-none">UID</label>
        <button class="px-4 py-1 cursor-pointer hover:bg-gray-300/50 hover:shadow hover:-top-1 select-none">{currentUID}</button>
        <!--{#await GetUIDs() then uidList}-->
        <!--    {#each uidList as uid, i}-->
        <!--        <div class="px-5 py-2 cursor-pointer hover:bg-gray-300/50 hover:shadow hover:-top-1 select-none"-->
        <!--             on:click={()=>currentUID=uid}>{uid}</div>-->
        <!--    {/each}-->
        <!--{/await}-->
    </div>
</div>