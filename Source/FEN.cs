namespace Chess;

// class for reading and handling FEN Notation
public class FEN
{
    // input info
    string notation;

    public FEN(string notation)
    {
        this.notation = notation;
    }

    // get a chessboard object representing the board position
    public Chessboard Board()
    {
        Chessboard b   = new();
        string[] ranks = this.notation.Split(' ')[0].Split('/');

        int y = 0, x = 0;
        foreach (string rank in ranks)
        {
            foreach (char piece in rank)
            {
                if (char.IsNumber(piece))
                    x += piece - 48;
                else
                {
                    PieceType pieceType = (PieceType)piece;
                    Color     color     = Color.Black;

                    if (char.IsUpper(piece))
                    {
                        pieceType = (PieceType)char.ToLower(piece);
                        color     = Color.White;
                    }

                    b.SetSquare(new Square(pieceType, color), y, x);
                    x++;
                }
            }
            y++;
            x = 0;
        }

        return b;
    }

    // get the side to move
    public Color activeColor()
    {
        if (notation.Split(' ')[1] == "w")
            return Color.White;
        return Color.Black;
    }

    // get a string representing the castling ability for both sides
    //  upper = white, lower = black
    //  q = queenside, k = kingside
    public string castlingAbility()
    {
        return notation.Split(' ')[2];
    }

    // en passant target square
    // null if no en enpassant exists
    public string? enPassantSquare()
    {
        string square = notation.Split(' ')[3];
        if (square == "-")
            return null;
        return square;
    }

    // halfmoves since the last pawn move or capture
    //  used for the 50 move rule
    public int halfmoveClock()
    {
        return Convert.ToInt32(notation.Split(' ')[4]);
    }

    // number of full moves in the game
    public int fullmoveClock()
    {
        return Convert.ToInt32(notation.Split(' ')[5]);
    }

    public string fen()
    {
        return this.notation;
    }

    public override string ToString()
    {
        return fen();
    }
}