export namespace api {
    export type WishType = number
    export type SharedWishType = number
    export type Rarity = number

    export const BeginnersWish: WishType = 100 // Beginners' Wish (Novice Wish)
    export const StandardWish: WishType = 200 // Standard Wish (Permanent Wish)
    export const CharacterEventWish: WishType = 301 // Character Event Wish
    export const WeaponEventWish: WishType = 302 // Weapon Event Wish
    export const CharacterEventWish2: WishType = 400 // Character Event Wish-2

    export const SBeginnersWish: SharedWishType = 100 // Beginners' Wish (Novice Wish)
    export const SStandardWish: SharedWishType = 200 // Standard Wish (Permanent Wish)
    export const SCharacterEventWish: SharedWishType = 301 // Character Event Wish and Character Event Wish-2
    export const SWeaponEventWish: SharedWishType = 302 // Weapon Event Wish

    export const Star1: Rarity = 1
    export const Star2: Rarity = 2
    export const Star3: Rarity = 3
    export const Star4: Rarity = 4
    export const Star5: Rarity = 5

    export const SharedWishTypes: SharedWishType[] = [SCharacterEventWish, SWeaponEventWish, SStandardWish, SBeginnersWish]
}