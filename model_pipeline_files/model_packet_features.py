import torch
import torch.nn as nn
import torch.nn.functional as F
# 1d


class GAB(nn.Module):
    def __init__(self, channels):
        super(GAB, self).__init__()
        # self.conv1 = nn.Conv2d(channels*16, channels*8, kernel_size=1) # 16
        self.conv2 = nn.Conv2d(channels, channels, kernel_size=1)
        self.conv3 = nn.Conv2d(channels, channels, kernel_size=1) #cghannels to be changed a bit
        self.global_avg_pool = nn.AdaptiveAvgPool2d(1)

    def forward(self, inputs):
        # x = F.relu(self.conv1(inputs))
        gmp = self.global_avg_pool(inputs)
        gmp = F.relu(self.conv2(gmp))
        gmp = F.relu(self.conv3(gmp))
        C_A = gmp * inputs
        x = torch.mean(C_A, dim=-3, keepdim=True)
        x = F.relu(x)
        G_OUT = x * C_A
        return G_OUT


class CAB(nn.Module):
    def __init__(self, classes, feature_per_class,channels):
        super(CAB, self).__init__()
        self.classes = classes
        self.feature_per_class = feature_per_class
        self.conv1 = nn.Conv2d(in_channels=channels, out_channels=feature_per_class * classes, kernel_size=1) #8
        # self.bn = nn.BatchNorm2d(feature_per_class * classes)
        self.dropout = nn.Dropout(0.5)
        self.global_max_pool = nn.AdaptiveMaxPool2d(1)

    def forward(self, inputs):
        F1 = self.conv1(inputs)
        F1 = F.relu(F1)
        F2 = self.dropout(F1)
        x = self.global_max_pool(F2)
        x = x.view(-1, self.classes, self.feature_per_class)
        S = torch.mean(x, dim=-1)
        # S = S.view(-1, self.classes)
        F_avg = torch.mean(F1.view(-1, self.classes, self.feature_per_class, inputs.shape[-2], inputs.shape[-1],), dim=-3)
        # print(S.shape,F_avg.shape)
        attention = S.unsqueeze(-1) * F_avg #.unsqueeze
        C_ATTN = torch.mean(attention, dim=-3, keepdim=True)
        # print(C_ATTN.shape,inputs.shape)
        C_OUT = inputs * C_ATTN
        return C_OUT


class Attention(nn.Module):
    def __init__(self, num_classes, feature_per_class,channels):
        super(Attention, self).__init__()
        self.GAB = GAB(channels)
        self.CAB = CAB(num_classes, feature_per_class,channels)
        self.global_avg_pool = nn.AdaptiveAvgPool2d(1)
        self.fc1 = nn.Linear(channels, 8) #num_classes * feature_per_class
        self.fc2 = nn.Linear(8, num_classes)

    def forward(self, inputs):
        x = self.GAB(inputs)
        x = self.CAB(x)
        x = self.global_avg_pool(inputs)
        x = x.view(x.size(0), -1)
        x = F.relu(self.fc1(x))
        x = self.fc2(x)
        return F.sigmoid(x)


class Hex1DModel(nn.Module):
    def __init__(self, num_classes, feature_per_class, input_shape):
        super(Hex1DModel, self).__init__()
        self.pad1 = nn.ZeroPad2d(1)
        self.pad2 = nn.ZeroPad2d(2)
        self.conv1 = nn.Conv2d(in_channels=input_shape[1], out_channels=8, kernel_size=3, padding=0)
        self.bn1 = nn.BatchNorm2d(8)
        self.conv2 = nn.Conv2d(in_channels=8, out_channels=16, kernel_size=5, padding=0)
        self.bn2 = nn.BatchNorm2d(16)
        self.maxpool1 = nn.MaxPool2d(kernel_size=2)
        self.maxpool2 = nn.MaxPool2d(kernel_size=2)
        self.attention = Attention(num_classes=num_classes, feature_per_class=feature_per_class,channels=16)

    def forward(self, inputs):
        x = self.pad1(inputs)
        x = F.relu(self.conv1(x))
        x = self.maxpool1(x)
        x = self.bn1(x)
        x = self.pad2(x)
        x = F.relu(self.conv2(x))
        x = self.bn2(x)
        x = self.maxpool2(x)
        x = self.attention(x)
        return x

if __name__=="__main__":
    args=parser.parse_args()
    obj=()
    obj.make_png_dir(args.input_dir)