interface AvatarProps {
    icon: string;
    name: string;
    pulls: number;
}

function Avatar({icon, name, pulls}: AvatarProps) {
    return (
        <div
            className="cursor-pointer bg-white/50 space-y-1.5 pt-4 flex flex-col w-fit items-center select-none  border shadow-sm rounded-lg hover:bg-white/25 transition duration-200">
            <div className="px-4">
                <div className="w-24">
                    <img alt={name} className="pointer-events-none rounded-full shadow-inner bg-amber-600"
                         src={icon}/>
                </div>
            </div>
            <div className="font-semibold text-gray-900 tracking-wide leading-relaxed">{name}</div>
            <div
                className="text-sm text-gray-700 border-t py-0.5 w-full text-center bg-gray-300/25 rounded-b-lg leading-relaxed tracking-wider">{pulls}</div>
        </div>
    )
}

export default Avatar