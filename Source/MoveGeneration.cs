using System.Windows.Markup;

namespace Chess;

public class MoveGenerator
{
    Chessboard board;

    public MoveGenerator(Chessboard board)
    {
        this.board = board;
    }

    // psuedo legal moves
    public List<Move> Moves()
    {
        List<Move> moves = new();

        for (int y = 0; y < 8; y++)
        for (int x = 0; x < 8; x++)
        {
            PieceType piece = board.GetSquare(x, y);
            if (Piece.GetPieceColor(piece) != board.activeColor())
                continue;
            
            switch (piece)
            {
                case PieceType.None:
                case PieceType.Invalid:
                    break;

                case PieceType.BlackRook:
                case PieceType.WhiteRook:
                    moves.AddRange(RookMoves(new Point(x, y)));
                    break;

                case PieceType.BlackKnight:
                case PieceType.WhiteKnight:
                    moves.AddRange(KnightMoves(new Point(x, y)));
                    break;

                case PieceType.BlackBishop:
                case PieceType.WhiteBishop:
                    moves.AddRange(BishopMoves(new Point(x, y)));
                    break;

                case PieceType.BlackQueen:
                case PieceType.WhiteQueen:
                    moves.AddRange(QueenMoves(new Point(x, y)));
                    break;

                case PieceType.BlackKing:
                case PieceType.WhiteKing:
                    moves.AddRange(KingMoves(new Point(x, y)));
                    break;

                case PieceType.BlackPawn:
                case PieceType.WhitePawn:
                    moves.AddRange(PawnMoves(new Point(x, y)));
                    break;
            }

        }

        return moves;
    }

    // moves of peices that slide
    public List<Move> SliderMoves(Point start, Point pattern)
    {
        List<Move> moves = new();

        int x = start.x + pattern.x;
        int y = start.y + pattern.y;

        while(this.board.GetSquare(x, y) == PieceType.None)
        {
            moves.Add(new Move(start, new Point(x, y)));
            x += pattern.x;
            y += pattern.y;
        }

        bool lastIsCapture  = this.board.GetSquare(x, y)                      != PieceType.Invalid;
        bool lastIsOppColor = Piece.GetPieceColor(this.board.GetSquare(x, y)) != board.activeColor();

        if (lastIsCapture && lastIsOppColor)
            moves.Add(new Move(start, new Point(x, y)));

        return moves;
    }

    public List<Move> BishopMoves(Point start)
    {
        List<Move> moves = new();

        moves.AddRange(this.SliderMoves(start, new Point(1, 1)));
        moves.AddRange(this.SliderMoves(start, new Point(-1, 1)));
        moves.AddRange(this.SliderMoves(start, new Point(1, -1)));
        moves.AddRange(this.SliderMoves(start, new Point(-1, -1)));

        return moves;
    }

    public List<Move> RookMoves(Point start)
    {
        List<Move> moves = new();

        moves.AddRange(this.SliderMoves(start, new Point(1, 0)));
        moves.AddRange(this.SliderMoves(start, new Point(0, 1)));
        moves.AddRange(this.SliderMoves(start, new Point(-1, 0)));
        moves.AddRange(this.SliderMoves(start, new Point(0, -1)));

        return moves;
    }

    public List<Move> QueenMoves(Point start)
    {
        List<Move> moves = new();

        moves.AddRange(this.RookMoves(start));
        moves.AddRange(this.BishopMoves(start));

        return moves;
    }

    public List<Move> KnightMoves(Point start)
    {
        List<Move>  moves  = new();
        List<Point> points = new();

        // possible moves
        points.Add(new Point(start.x + 2, start.y + 1));
        points.Add(new Point(start.x + 1, start.y + 2));
        points.Add(new Point(start.x + 2, start.y - 1));
        points.Add(new Point(start.x + 1, start.y - 2));
        points.Add(new Point(start.x - 1, start.y - 2));
        points.Add(new Point(start.x - 2, start.y - 1));
        points.Add(new Point(start.x - 2, start.y + 1));
        points.Add(new Point(start.x - 1, start.y + 2));

        // psuedo-legal moves
        foreach (Point point in points)
        {
            PieceType square = this.board.GetSquare(point);

            // true if the piece we're capturing is opposite color to the one from the origin point
            bool capIsOppColor = 
               Piece.GetPieceColor(this.board.GetSquare(start))
            != board.activeColor();

            // add if valid move or capture
            if (square == PieceType.None || (square != PieceType.Invalid && capIsOppColor))
                moves.Add(new Move(start, point));
        }

        return moves;
    }

    public List<Move> KingMoves(Point start)
    {
        List<Move>  moves  = new();
        List<Point> points = new();

        // possible moves
        points.Add(new Point(start.x + 1, start.y + 1));
        points.Add(new Point(start.x + 1, start.y - 1));
        points.Add(new Point(start.x - 1, start.y + 1));
        points.Add(new Point(start.x - 1, start.y - 1));
        points.Add(new Point(start.x + 1, start.y + 0));
        points.Add(new Point(start.x + 0, start.y + 1));
        points.Add(new Point(start.x - 1, start.y + 0));
        points.Add(new Point(start.x + 0, start.y - 1));

        // psuedo-legal moves
        foreach (Point point in points)
        {
            PieceType square = this.board.GetSquare(point);

            // true if the piece we're capturing is opposite color
            bool capIsOppColor = 
               Piece.GetPieceColor(this.board.GetSquare(start))
            != board.activeColor();

            // add if valid move or capture
            if (square == PieceType.None || (square != PieceType.Invalid && capIsOppColor))
                moves.Add(new Move(start, point));
        }

        return moves;
    }

    public List<Move> PawnMoves(Point start)
    {
        List<Move> moves = new();

        int mod = 1;
        if (board.activeColor() == Color.Black)
            mod = -1;

        // moves
        Point p1 = new Point(start.x, start.y + (mod * 1));
        if (board.GetSquare(p1) == PieceType.None)
        {
            moves.Add(new Move(start, p1));
            Point p2 = new Point(start.x, start.y + (mod * 2));
            if (board.GetSquare(p2) == PieceType.None)
                if (board.activeColor() == Color.White && start.y == 1)
                    moves.Add(new Move(start, p2));
                else if (board.activeColor() == Color.Black && start.y == 7)
                    moves.Add(new Move(start, p2));
        }

        // captures
        Point p3 = new Point(start.x + 1, start.y + 1);
        bool isCapture   = board.GetSquare(p3) != PieceType.None && board.GetSquare(p3) != PieceType.Invalid;
        bool isOppColor  = Piece.GetPieceColor(board.GetSquare(p3)) != board.activeColor();
        
        bool isEnPassant = false;
        if (board.enPassantSquare() != null)
            isEnPassant = board.enPassantSquare().x == p3.x && board.enPassantSquare().y == p3.y;

        if ((isCapture && isOppColor) || isEnPassant)
            moves.Add(new Move(start, p3));

        Point p4 = new Point(start.x - 1, start.y + 1);
        isCapture   = board.GetSquare(p4) != PieceType.None && board.GetSquare(p4) != PieceType.Invalid;
        isOppColor  = Piece.GetPieceColor(board.GetSquare(p4)) != board.activeColor();
        isEnPassant = false;
        if (board.enPassantSquare() != null)
            isEnPassant = board.enPassantSquare().x == p4.x && board.enPassantSquare().y == p4.y;

        if ((isCapture && isOppColor) || isEnPassant)
            moves.Add(new Move(start, p4));
        
        return moves;
    }
}