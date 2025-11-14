export default function ColorTest() {
  return (
    <div className="p-8 space-y-4">
      <h1 className="text-3xl font-bold">Paleta CarPooling</h1>

      <div className="flex gap-4">
        <div className="w-32 h-32 bg-primary rounded-lg flex items-center justify-center text-white font-bold">
          Primary
        </div>
        <div className="w-32 h-32 bg-secondary rounded-lg flex items-center justify-center text-white font-bold">
          Secondary
        </div>
        <div className="w-32 h-32 bg-accent rounded-lg flex items-center justify-center text-white font-bold">
          Accent
        </div>
        <div className="w-32 h-32 bg-success rounded-lg flex items-center justify-center text-white font-bold">
          Success
        </div>
      </div>

      <div className="space-y-2">
        <button className="px-6 py-3 bg-primary text-white rounded-lg hover:bg-primary-600 transition">
          Botón Primary
        </button>
        <button className="px-6 py-3 bg-secondary text-white rounded-lg hover:bg-secondary-600 transition">
          Botón Secondary
        </button>
        <button className="px-6 py-3 bg-accent text-white rounded-lg hover:bg-accent/90 transition">
          Botón Accent
        </button>
      </div>
    </div>
  )
}
