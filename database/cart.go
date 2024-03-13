package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("cant find product")
	ErrCantDecodeProduct  = errors.New("cant find the product(decode)")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cant update ")
	ErrCantRemoveItemCart = errors.New("cant remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
