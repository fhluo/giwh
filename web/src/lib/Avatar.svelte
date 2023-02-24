<script lang="ts">
    export let icon: string
    export let name: string
    export let pulls: number
</script>

<div class="cursor-pointer bg-white/50 space-y-1.5 pt-4 flex flex-col w-fit items-center select-none  border shadow-sm rounded-lg hover:bg-white/25 transition duration-200">
    <div class="px-4 ">
        <div class="w-24">
            <img alt="{name}" class="pointer-events-none rounded-full shadow-inner bg-amber-600"
                 src="{icon}">
        </div>
    </div>
    <div class="font-semibold text-gray-900 tracking-wide leading-relaxed">{name}</div>
    <div class="text-sm text-gray-700 border-t py-0.5 w-full text-center bg-gray-300/25 rounded-b-lg leading-relaxed tracking-wider">{pulls}</div>
</div>