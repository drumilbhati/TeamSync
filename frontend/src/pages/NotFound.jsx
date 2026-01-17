import { Link } from "react-router-dom";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

const NotFound = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] bg-background text-foreground text-center px-4">
      <AlertCircle className="w-24 h-24 text-destructive mb-6" />
      <h1 className="text-4xl font-bold tracking-tight mb-2">404</h1>
      <h2 className="text-2xl font-semibold mb-4">Page Not Found</h2>
      <p className="text-muted-foreground max-w-md mb-8">
        Sorry, we couldn't find the page you're looking for. It might have been
        removed, renamed, or doesn't exist.
      </p>
      <Link to="/">
        <Button size="lg">
          Go back home
        </Button>
      </Link>
    </div>
  );
};

export default NotFound;
