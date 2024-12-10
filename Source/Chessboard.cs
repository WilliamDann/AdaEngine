using System.Text;

namespace Chess;

// Square-Centric representation of a chess board
public class Chessboard
{
    Square[,] board;

    public Chessboard()
    {
        this.board = new Square[12, 12];
        this.Clear();
    }

    // get what is at a given square
    public Square GetSquare(int x, int y)
    {
        // offset past buffer
        x += 2;
        y += 2;

        return this.board[x, y];
    }

    // plcae a piece on a square
    public void SetSquare(Square piece, int x, int y)
    {
        // offset for buffer
        x += 2;
        y += 2;

        // bounds checking
        if (x < 2 || x > 10 || y < 2 || y > 10)
            throw new IndexOutOfRangeException($"(${x},{y}) is an invalid square to place a piece");

        // place
        this.board[x, y] = piece;
    }

    // clear a squrae
    public void ClearSquare(int x, int y)
    {
        // offset for buffer
        x += 2;
        y += 2;

        // determine clear square color
        Color c = Color.Black;
        if (x % 2 == 0 ^ y % 2 == 0)
            c = Color.White;

        // replace square
        this.board[x, y] = new Square(PieceType.None, c);
    }

    // clear the board to it's default state
    // Invalid squares for 2 outer layers, None squares for inner ones
    public void Clear()
    {
        for (int x = 0; x < 12; x++)
            for (int y = 0; y < 12; y++)
            {
                // set the color correctly for an empty square
                Color c = Color.Black;
                if (x % 2 == 0 ^ y % 2 == 0)
                    c = Color.White;

                if (x <= 1 || x >= 10 || y <= 1 || y >= 10)
                    this.board[x, y] = new Square(PieceType.Invalid, c);
                else
                    this.board[x, y] = new Square(PieceType.None, c);
            }
    }

    // Iter over valid board squares
    public IEnumerable<Square> SquareIter()
    {
        for (int x = 2; x < 10; x++)
            for (int y = 2; y < 10; y++)
                yield return this.board[x, y];
    }

    // print the board to a string
    public override string ToString()
    {
        StringBuilder sb = new StringBuilder();
        for (int x = 2; x < 10; x++)
        {
            for (int y = 2; y < 10; y++)
            {
                sb.Append(this.board[x, y].Char());
                sb.Append(" ");
            }
            sb.Append("\n");
        }
        return sb.ToString();
    }
}