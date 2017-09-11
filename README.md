# tripoli
Package tripoli implements a simplified version of the Michigan portion of the card game Tripoli.

## Rules
The game is played with a standard 52-card deck. Aces are high. A dealer is chosen at random for the first hand, and the deal rotates to the left on each subsequent hand.

Before each hand, players ante one chip into each of five stake pots labeled Ten, Jack, Queen, King, and Ace of Hearts, called the “counter” cards, and an additional pot called the Kitty. Starting with the player to the dealer's left, the entire deck is dealt out into a hand for each player plus an “extra hand” after the dealer's hand, which remains face down and is not played.

The player holding the lowest club begins play by discarding it. When a card is discarded, whoever holds the next higher card in the same suit must discard it, and so on. When an ace is played or no player holds the next card, whoever played the last card must restart play with their lowest card in either of the suits of the opposite color. A player who plays a counter collects all of the chips in the corresponding pot.

If a player is unable to restart play because they do not hold any cards of the required color, they must pay one chip to the Kitty and control of the restart passes to the player to their left. If no player holds any cards of the required color, then after every player has paid one chip to the Kitty consecutively, the hand is over. Otherwise, the hand is won by the player who plays their last card. When the hand is over, every player must pay one chip to the Kitty for each card remaining in their hand. Then the winner, if there is one, collects the Kitty. Any unclaimed stakes remain on the table for the following hand.
