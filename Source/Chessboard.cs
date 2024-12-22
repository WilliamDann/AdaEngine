using System.Text;

namespace Chess;

// Square-Centric representation of a chess board
public class Chessboard : FEN
{
    PieceType[,] board;

    public Chessboard(string initalPosition = INITAL) : base(initalPosition)
    {
        this.board = new PieceType[12, 12];
        this.Update();
    }

    // update the board from the FEN position
    public void Update()
    {
        this.Clear();

        string[] ranks = notation.Split(' ')[0].Split('/');

        int y = 7, x = 7;
        foreach (string rank in ranks)
        {
            foreach (char piece in rank)
                if (char.IsNumber(piece))
                    x -= piece - 48;
                else
                    SetSquare((PieceType)piece, new Point(x--, y));

            y--;
            x = 7;
        }
    }

    // get what is at a given square
    public PieceType GetSquare(Point square)
    {
        return this.board[square.y + 2, square.x + 2];
    }

    public PieceType GetSquare(int x, int y)
    {
        return this.board[y + 2, x + 2];
    }

    // plcae a piece on a square
    public void SetSquare(PieceType piece, Point square)
    {
        Point p = new(square.x + 2, square.y + 2);

        // bounds checking
        if (p.x < 2 || p.x > 10 || p.y < 2 || p.y > 10)
            throw new IndexOutOfRangeException($"(${square.x},{square.y}) is an invalid square to place a piece");

        // place
        this.board[p.y, p.x] = piece;
    }

    // clear a squrae
    public void ClearSquare(Point square)
    {
        // replace square
        this.board[square.y + 2, square.x + 2] = PieceType.None;
    }

    // clear the board to it's default state
    // Invalid squares for 2 outer layers, None squares for inner ones
    public void Clear()
    {
        for (int x = 0; x < 12; x++)
            for (int y = 0; y < 12; y++)
                if (x <= 1 || x >= 10 || y <= 1 || y >= 10)
                    this.board[x, y] = PieceType.Invalid;
                else
                    this.board[x, y] = PieceType.None;
    }

    // Iter over valid board squares
    public IEnumerable<PieceType> SquareIter()
    {
        for (int x = 2; x < 10; x++)
            for (int y = 2; y < 10; y++)
                yield return this.board[y, x];
    }

    // print the board to a string
    public override string ToString()
    {
        StringBuilder sb = new StringBuilder();
        for (int x = 2; x < 10; x++)
        {
            for (int y = 2; y < 10; y++)
                sb.Append((char)this.board[x, y] + " ");
            sb.Append("\n");
        }
        return sb.ToString();
    }
}