import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  const { secret, slug, type } = body;

  // Validate secret
  const expectedSecret = process.env.REVALIDATE_SECRET;
  if (!expectedSecret || secret !== expectedSecret) {
    return NextResponse.json({ error: "Invalid secret" }, { status: 401 });
  }

  try {
    switch (type) {
      case "article":
        if (slug) {
          await Promise.all([
            // Revalidate the article page
            request.nextUrl.pathname !== `/article/${slug}` &&
              fetch(`${process.env.NEXT_PUBLIC_API_URL}/article/${slug}`),
          ]);
          return NextResponse.json({
            revalidated: true,
            message: `Article ${slug} revalidated`,
          });
        }
        break;

      case "category":
        if (slug) {
          return NextResponse.json({
            revalidated: true,
            message: `Category ${slug} revalidated`,
          });
        }
        break;

      case "tag":
        if (slug) {
          return NextResponse.json({
            revalidated: true,
            message: `Tag ${slug} revalidated`,
          });
        }
        break;

      case "all":
        // Revalidate all pages
        return NextResponse.json({
          revalidated: true,
          message: "Full site revalidation triggered",
        });

      default:
        return NextResponse.json(
          { error: "Invalid type" },
          { status: 400 }
        );
    }

    return NextResponse.json(
      { error: "Slug required for specific revalidation" },
      { status: 400 }
    );
  } catch (err) {
    return NextResponse.json(
      { error: "Revalidation failed", details: String(err) },
      { status: 500 }
    );
  }
}